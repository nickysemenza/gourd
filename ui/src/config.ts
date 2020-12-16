import Cookies from "universal-cookie";
import { Configuration } from "./api/openapi-fetch";
import jwt_decode from "jwt-decode";
import { toast } from "react-toastify";

export const getAPIURL = () => getBaseURL() + "/api";
export const getGQLURL = () => getBaseURL() + "/query";
export const getBaseURL = () => process.env.REACT_APP_API_URL;

export const getConfig = () => {
  // TODO: load this all from API + cache it
  return {
    google: {
      scopes:
        "profile email https://www.googleapis.com/auth/photoslibrary.readonly",
      client_id:
        "520431142247-bcog816pfdh6bctlvbreh3i3urhpidv5.apps.googleusercontent.com",
    },
  };
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
  toast.error(err.message);
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
