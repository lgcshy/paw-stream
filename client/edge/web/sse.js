/**
 * SSE (Server-Sent Events) Client
 */
class SSEClient {
    constructor(url) {
        this.url = url;
        this.eventSource = null;
        this.listeners = {
            connected: [],
            status: [],
            log: [],
            config: [],
            error: [],
            close: []
        };
        this.reconnectDelay = 3000;
        this.maxReconnectDelay = 30000;
        this.reconnectAttempts = 0;
    }

    /**
     * Connect to SSE endpoint
     */
    connect() {
        if (this.eventSource) {
            this.eventSource.close();
        }

        console.log('[SSE] Connecting to', this.url);
        this.eventSource = new EventSource(this.url);

        this.eventSource.onopen = () => {
            console.log('[SSE] Connected');
            this.reconnectAttempts = 0;
            this.reconnectDelay = 3000;
            this.trigger('connected');
        };

        this.eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                // For events without explicit type, treat as generic message
                this.handleEvent({
                    type: event.type || 'message',
                    data: data
                });
            } catch (err) {
                console.error('[SSE] Failed to parse event:', err, event.data);
            }
        };

        // Listen for specific event types
        this.eventSource.addEventListener('connected', (event) => {
            try {
                const data = JSON.parse(event.data);
                console.log('[SSE] Server acknowledged connection:', data);
            } catch (err) {
                console.error('[SSE] Failed to parse connected event:', err);
            }
        });

        this.eventSource.addEventListener('status', (event) => {
            try {
                const data = JSON.parse(event.data);
                this.trigger('status', data);
            } catch (err) {
                console.error('[SSE] Failed to parse status event:', err);
            }
        });

        this.eventSource.addEventListener('log', (event) => {
            try {
                const data = JSON.parse(event.data);
                this.trigger('log', data);
            } catch (err) {
                console.error('[SSE] Failed to parse log event:', err);
            }
        });

        this.eventSource.addEventListener('config', (event) => {
            try {
                const data = JSON.parse(event.data);
                this.trigger('config', data);
            } catch (err) {
                console.error('[SSE] Failed to parse config event:', err);
            }
        });

        this.eventSource.onerror = (error) => {
            console.error('[SSE] Error:', error);
            this.trigger('error', error);
            
            if (this.eventSource.readyState === EventSource.CLOSED) {
                this.handleReconnect();
            }
        };
    }

    /**
     * Handle SSE event
     */
    handleEvent(data) {
        const { type, data: payload } = data;
        
        switch (type) {
            case 'status':
                this.trigger('status', payload);
                break;
            case 'log':
                this.trigger('log', payload);
                break;
            case 'config':
                this.trigger('config', payload);
                break;
            case 'connected':
                console.log('[SSE] Server acknowledged connection:', payload);
                break;
            default:
                console.warn('[SSE] Unknown event type:', type, payload);
        }
    }

    /**
     * Handle reconnection
     */
    handleReconnect() {
        this.reconnectAttempts++;
        const delay = Math.min(
            this.reconnectDelay * Math.pow(1.5, this.reconnectAttempts - 1),
            this.maxReconnectDelay
        );

        console.log(`[SSE] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})...`);
        
        setTimeout(() => {
            this.connect();
        }, delay);
    }

    /**
     * Register event listener
     */
    on(event, callback) {
        if (this.listeners[event]) {
            this.listeners[event].push(callback);
        } else {
            console.warn('[SSE] Unknown event type:', event);
        }
    }

    /**
     * Trigger event
     */
    trigger(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(callback => {
                try {
                    callback(data);
                } catch (err) {
                    console.error('[SSE] Listener error:', err);
                }
            });
        }
    }

    /**
     * Disconnect
     */
    disconnect() {
        if (this.eventSource) {
            console.log('[SSE] Disconnecting');
            this.eventSource.close();
            this.eventSource = null;
            this.trigger('close');
        }
    }

    /**
     * Get connection state
     */
    getState() {
        if (!this.eventSource) {
            return 'disconnected';
        }

        switch (this.eventSource.readyState) {
            case EventSource.CONNECTING:
                return 'connecting';
            case EventSource.OPEN:
                return 'connected';
            case EventSource.CLOSED:
                return 'disconnected';
            default:
                return 'unknown';
        }
    }
}
