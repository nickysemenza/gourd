import React, { useState } from "react";
import { useCookies } from "react-cookie";

import {
  GoogleLoginResponse,
  GoogleLoginResponseOffline,
  useGoogleLogin,
} from "react-google-login";
import { AuthResp, Configuration, DefaultApi } from "../api/openapi-fetch";
import { getAPIURL, getConfig } from "../config";
import Debug from "./Debug";

const Login: React.FC = () => {
  const api = new DefaultApi(new Configuration({ basePath: getAPIURL() }));
  const cookieName = "jwt";
  const [cookies, setCookie, removeCookie] = useCookies(["cookie-name"]);

  const [auth, setAuth] = useState<AuthResp>();

  const onSuccess = async (
    response: GoogleLoginResponse | GoogleLoginResponseOffline
  ) => {
    console.log({ response });
    const { code } = response;
    if (code !== undefined) {
      const resp = await api.authLogin({ code });
      console.log({ resp });
      setCookie(cookieName, resp.jwt);
      setAuth(resp);
    } else {
      throw new Error("bad auth" + response);
    }
  };
  const { google } = getConfig();
  const { signIn, loaded } = useGoogleLogin({
    onSuccess,
    // onAutoLoadFinished,
    clientId: google.client_id,
    // cookiePolicy,
    // loginHint,
    // hostedDomain,
    // autoLoad,
    // isSignedIn,
    // fetchBasicProfile,
    // redirectUri,
    // discoveryDocs,
    // onFailure,
    // uxMode,
    scope: google.scopes,
    accessType: "offline",
    responseType: "code",
    // jsSrc,
    // onRequest,
    // prompt,
  });

  return (
    <div>
      {loaded && <button onClick={signIn}>login</button>}
      <Debug data={{ auth, cookies }} />
    </div>
  );
};

export default Login;
