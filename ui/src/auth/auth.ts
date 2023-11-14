import Cookies from "universal-cookie";
import { jwtDecode } from "jwt-decode";

export const COOKIE_NAME = "gourd-jwt";

export const getJWT = (): string | undefined => {
  const cookies = new Cookies();
  return cookies.get(COOKIE_NAME);
};

export const parseJWT = (): JWT | undefined => {
  const jwt = getJWT();
  if (jwt === "" || !jwt) return;
  const d: JWT = jwtDecode(jwt);
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
