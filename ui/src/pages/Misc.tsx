import React from "react";
import { getBaseURL } from "../config";
import { RedocStandalone } from "redoc";

export const Docs: React.FC = () => {
  return <RedocStandalone specUrl={getBaseURL() + "/spec"} />;
};
