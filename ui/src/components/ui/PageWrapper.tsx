import { Helmet } from "react-helmet";

const PageWrapper: React.FC<{
  children: React.ReactNode;
  title: string;
}> = ({ children, title }) => {
  return (
    <>
      {/* @ts-ignore */}
      <Helmet>
        <title>{title} | gourd</title>
      </Helmet>
      {children}
    </>
  );
};

export default PageWrapper;
