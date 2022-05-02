import Cookies from "universal-cookie";
import { Configuration, SystemApi } from "./api/openapi-fetch";
import jwt_decode from "jwt-decode";
import { toast } from "react-toastify";
import { json } from "stream/consumers";

export const getAPIURL = () => getBaseURL() + "/api";
export const getBaseURL = () => process.env.REACT_APP_API_URL;
export const getTracingURL = () => process.env.REACT_APP_TRACING_URL || "";

export const getConfig = async () => {
  const a = new SystemApi(getOpenapiFetchConfig());
  let config = await a.getConfig();
  return config;
};

export const COOKIE_NAME = "gourd-jwt";

export const getJWT = (): string | undefined => {
  const cookies = new Cookies();
  return cookies.get(COOKIE_NAME);
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

export const parseJWT = (): JWT | undefined => {
  const jwt = getJWT();
  if (jwt === "" || !jwt) return;
  const d: JWT = jwt_decode(jwt);
  return d;
};
export const isLoggedIn = () => {
  const t = parseJWT();
  if (!t || t.exp < Math.floor(Date.now() / 1000)) return false;
  return true;
};

export const getName = () => {
  if (!isLoggedIn()) return;
  const t = parseJWT();
  return t?.user_info.name;
};

export interface UserInfo {
  email: string;
  family_name: string;
  given_name: string;
  id: string;
  locale: string;
  name: string;
  picture: string;
  verified_email: boolean;
}

export interface JWT {
  user_info: UserInfo;
  exp: number;
  iat: number;
}
