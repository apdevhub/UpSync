import React from 'react';

// ─── Props ────────────────────────────────────────────────────
interface ErrorBoundaryState {
    hasError: boolean;
    error: Error | null;
}

interface ErrorBoundaryProps {
    children: React.ReactNode;
    fallback?: React.ReactNode;
}

// ─── Global Error Boundary ────────────────────────────────────
export class ErrorBoundary extends React.Component<
    ErrorBoundaryProps,
    ErrorBoundaryState
> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = { hasError: false, error: null };
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return { hasError: true, error };
    }

    componentDidCatch(error: Error, info: React.ErrorInfo) {
        console.error('[ErrorBoundary] Uncaught error:', error, info);
    }

    handleReset = () => {
        this.setState({ hasError: false, error: null });
    };

    render() {
        if (this.state.hasError) {
            if (this.props.fallback) return this.props.fallback;

            return (
                <div className="share-page">
                    <div className="status-card">
                        <span className="status-icon">💥</span>
                        <h2>Something went wrong</h2>
                        <p>
                            {this.state.error?.message ?? 'An unexpected error occurred.'}
                            <br />
                            Please refresh the page or try again.
                        </p>
                        <button
                            className="go-home"
                            style={{ marginTop: 20, border: 'none', cursor: 'pointer' }}
                            onClick={this.handleReset}
                        >
                            🔄 Try Again
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
