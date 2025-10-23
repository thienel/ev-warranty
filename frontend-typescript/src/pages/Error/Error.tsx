import React from "react";

interface ErrorProps {
  code: number;
}

const Error: React.FC<ErrorProps> = ({ code }) => {
  return <div>Error {code}</div>;
};

export default Error;
