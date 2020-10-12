import React, { useState } from "react";
import { useCookies } from "react-cookie";

import {
  GoogleLoginResponse,
  GoogleLoginResponseOffline,
  useGoogleLogin,
} from "react-google-login";
import {
  AuthResp,
  Configuration,
  AuthenticationApi,
} from "../api/openapi-fetch";
import {
  COOKIE_NAME,
  getAPIURL,
  getConfig,
  getName,
  isLoggedIn,
} from "../config";
import Debug from "./Debug";

const Login: React.FC = () => {
  const api = new AuthenticationApi(
    new Configuration({ basePath: getAPIURL() })
  );
  const [cookies, setCookie] = useCookies(["cookie-name"]);

  const [auth, setAuth] = useState<AuthResp>();

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
    onFailure,
    // uxMode,
    scope: google.scopes,
    accessType: "offline",
    responseType: "code",
    // jsSrc,
    // onRequest,
    // prompt,
  });

  const loggedIn = isLoggedIn();
  return (
    <div>
      {!loggedIn ? (
        loaded && <button onClick={signIn}>login</button>
      ) : (
        <button>logged in as {getName()}</button>
      )}
      <Debug data={{ auth, cookies }} />
    </div>
  );
};

export default Login;
