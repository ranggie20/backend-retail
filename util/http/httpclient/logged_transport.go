package httpclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/online-bnsp/backend/util/logger"
	"github.com/rs/zerolog"
)

type LoggedTransport struct {
	zLogger       zerolog.Logger
	rtt           http.RoundTripper
	logModuleName string
}

func NewLoggedTransport(zLogger zerolog.Logger, rtt http.RoundTripper) *LoggedTransport {
	return &LoggedTransport{
		logModuleName: "http.client",
		zLogger:       zLogger,
		rtt:           rtt,
	}
}

func (t *LoggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var ct clientTrace
	req = req.WithContext(withTraceContext(req.Context(), &ct))
	start := time.Now()
	res, err := t.rtt.RoundTrip(req)
	duration := time.Since(start)
	ctxLogger := logger.FromContext(req.Context())
	if ctxLogger.GetLevel() == zerolog.Disabled {
		ctxLogger = &t.zLogger
	}
	lvl := ctxLogger.Info()
	if err != nil {
		lvl = ctxLogger.Error().Err(err)
	}
	var statusCode int
	if res != nil {
		statusCode = res.StatusCode
	}
	lvl.
		Str("module", t.logModuleName). // to make sure context or non-context logger has this value (for log filtering)
		Str("req_method", req.Method).
		Str("req_path", req.URL.Path).
		Str("req_query", req.URL.RawQuery).
		Str("req_scheme", req.URL.Scheme).
		Str("req_host", req.URL.Host).
		Int("res_status", statusCode).
		Int64("duration", duration.Milliseconds()).
		Int64("duration_0_dns_lookup", ct.dnsLookupTotal.Milliseconds()).
		Int64("duration_1_tcp_conn", ct.tcpConnectionTotal.Milliseconds()).
		Int64("duration_2_tls_handshake", ct.tlsHandshakeTotal.Milliseconds()).
		Int64("duration_3_server_process", ct.serverProcessingTotal.Milliseconds()).
		Bool("client_conn_reuse", ct.isReused).
		Bool("client_conn_idle", ct.isIdle).
		Int64("client_conn_idle_time", ct.idleDuration.Milliseconds()).
		Bool("client_dns_coalesced", ct.isDNSCoalesced).
		Bool("client_tls", ct.isTLS).
		Msgf("HTTP client %s %s://%s%s", req.Method, req.URL.Scheme, req.URL.Host, req.URL.Path)
	return res, err
}

// clientTrace copied from https://github.com/tcnksm/go-httpstat and adjusted accordingly.
type clientTrace struct {
	dnsStart    time.Time
	dnsDone     time.Time
	tcpStart    time.Time
	tcpDone     time.Time
	tlsStart    time.Time
	tlsDone     time.Time
	serverStart time.Time
	serverDone  time.Time
	// Connection phase breakdown
	dnsLookupTotal        time.Duration
	tcpConnectionTotal    time.Duration
	tlsHandshakeTotal     time.Duration
	serverProcessingTotal time.Duration
	// idle time for conn reuse = true
	idleDuration time.Duration
	// isTLS is true when connection seems to use TLS
	isTLS bool
	// isReused is true when connection is reused (keep-alive)
	isReused bool
	// isIdle is true if the connection was idle.
	isIdle bool
	// isDNSCoalesced whether  another caller who was doing the same DNS lookup concurrently.
	isDNSCoalesced bool
}

func withTraceContext(ctx context.Context, ct *clientTrace) context.Context {
	return httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			// if reuse, DNSStart(Done) and ConnectStart(Done) is skipped.
			ct.isReused = info.Reused
			ct.isIdle = info.WasIdle
			ct.idleDuration = info.IdleTime
		},
		GotFirstResponseByte: func() {
			ct.serverDone = time.Now()
			ct.serverProcessingTotal = ct.serverDone.Sub(ct.serverStart)
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			ct.dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			ct.dnsDone = time.Now()
			ct.dnsLookupTotal = ct.dnsDone.Sub(ct.dnsStart)
			ct.isDNSCoalesced = info.Coalesced
		},
		ConnectStart: func(network, addr string) {
			ct.tcpStart = time.Now()
			if ct.dnsStart.IsZero() {
				ct.dnsStart = ct.tcpStart
				ct.dnsDone = ct.tcpStart
			}
		},
		ConnectDone: func(network, addr string, err error) {
			ct.tcpDone = time.Now()
			ct.tcpConnectionTotal = ct.tcpDone.Sub(ct.tcpStart)
		},
		TLSHandshakeStart: func() {
			ct.isTLS = true
			ct.tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			ct.tlsDone = time.Now()
			ct.tlsHandshakeTotal = ct.tlsDone.Sub(ct.tlsStart)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			ct.serverStart = time.Now()
			// When connection is re-used, DNS/TCP/TLS hook is not called.
			if ct.isReused {
				now := ct.serverStart
				ct.dnsStart = now
				ct.dnsDone = now
				ct.tcpStart = now
				ct.tcpDone = now
				ct.tlsStart = now
				ct.tlsDone = now
			}
			if ct.isTLS {
				return
			}
			ct.tlsHandshakeTotal = ct.tcpDone.Sub(ct.tcpDone)
		},
	})
}
