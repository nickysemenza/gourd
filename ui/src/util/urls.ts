export const getBaseURL = () => import.meta.env.VITE_API_URL;
export const getTracingURL = () => import.meta.env.VITE_TRACING_URL || "";
export const getAPIURL = () => getBaseURL() + "/api";
