import React, { useState } from "react";
import { useCookies } from "react-cookie";

import {
  GoogleLoginResponse,
  GoogleLoginResponseOffline,
  useGoogleLogin,
} from "react-google-login";
import { AuthResp, AuthenticationApi, ConfigData } from "../api/openapi-fetch";
import { COOKIE_NAME, isLoggedIn, getName } from "../auth";
import { getOpenapiFetchConfig } from "../config";

const Login: React.FC<{ config: ConfigData }> = ({ config }) => {
  const api = new AuthenticationApi(getOpenapiFetchConfig());
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_cookies, setCookie] = useCookies([COOKIE_NAME]);
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_auth, setAuth] = useState<AuthResp>();

  const onSuccess = async (
    response: GoogleLoginResponse | GoogleLoginResponseOffline
  ) => {
    console.log({ response });
    const { code } = response;
    if (code !== undefined) {
      const resp = await api.authLogin({ code });
      console.log({ resp });
      setCookie(COOKIE_NAME, resp.jwt);
      setAuth(resp);
    } else {
      throw new Error("bad auth" + response);
    }
  };
  const onFailure = async (error: any) => {
    console.log("error", error);
  };

  const { signIn, loaded } = useGoogleLogin({
    onSuccess,
    // onAutoLoadFinished,
    clientId: config?.google_client_id || "",
    // cookiePolicy,
    // loginHint,
    // hostedDomain,
    // autoLoad,
    // isSignedIn,
    // fetchBasicProfile,
    // redirectUri,
    // discoveryDocs,
    onFailure,
    // uxMode,
    scope: config?.google_scopes || "",
    accessType: "offline",
    responseType: "code",
    // jsSrc,
    // onRequest,
    // prompt,
  });

  const loggedIn = isLoggedIn();
  return config ? (
    <div>
      {!loggedIn ? (
        loaded && <button onClick={signIn}>login</button>
      ) : (
        <button>logged in as {getName()}</button>
      )}
      {/* <Debug data={{ auth, cookies }} /> */}
    </div>
  ) : null;
};

export default Login;
