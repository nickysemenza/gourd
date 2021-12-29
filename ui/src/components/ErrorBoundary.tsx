import React, { Component, ErrorInfo, ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
  };

  public static getDerivedStateFromError(e: Error): State {
    // Update state so the next render will show the fallback UI.
    return { hasError: true, error: e };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Uncaught error:", error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div>
          <h1 className="text-xl">
            Sorry.. there was an error: {this.state.error?.name}
          </h1>
          {this.state.error?.message}
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
// https://react-typescript-cheatsheet.netlify.app/docs/basic/getting-started/error_boundaries/
