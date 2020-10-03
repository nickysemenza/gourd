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
