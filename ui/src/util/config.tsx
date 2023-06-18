import { Configuration, SystemApi } from "../api/openapi-fetch";
import { toast } from "react-toastify";
import { getJWT } from "../auth/auth";

export const getAPIURL = () => getBaseURL() + "/api";
export const getBaseURL = () => import.meta.env.VITE_API_URL;
export const getTracingURL = () => import.meta.env.VITE_TRACING_URL || "";

export const getConfig = async () => {
  const a = new SystemApi(getOpenapiFetchConfig());
  const config = await a.getConfig();
  return config;
};

export const getOpenapiFetchConfig = () => {
  const c = new Configuration({ basePath: getAPIURL(), accessToken: getJWT() });
  return c;
};

export const onAPIRequest = (req: Request): void => {
  req.headers.set("Authorization", "Bearer " + getJWT());
};
export const onAPIError = (
  err: {
    message: string;
    data: any;
    status?: number;
  },
  retry: () => Promise<any>,
  response?: Response
) => {
  console.log(err);
  const apiErr: Error = err.data;
  toast.error(
    <div>
      <div className="font-bold">{err.message}</div>
      <div>{apiErr.message}</div>
    </div>
  );
};
