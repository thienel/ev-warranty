import React, { useEffect } from "react";
import LoadingOverlay from "@components/LoadingOverlay/LoadingOverlay";
import { useNavigate, useSearchParams } from "react-router-dom";
import { loginSuccess, setToken, logout } from "@redux/authSlice";
import { message } from "antd";
import { useDispatch } from "react-redux";
import api from "@services/api";
import { API_ENDPOINTS } from "@constants/common-constants";
import useHandleApiError from "@/hooks/useHandleApiError";

const CallBack: React.FC = () => {
  const [searchParams] = useSearchParams();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const handleError = useHandleApiError();

  useEffect(() => {
    const handleLogin = async (): Promise<void> => {
      const token = searchParams.get("token");
      if (!token) {
        message.error("No token provided");
        navigate("/login");
        return;
      }

      try {
        // Set token first so it's available for the validation request
        dispatch(setToken(token));

        // Validate the token with the backend
        const res = await api.get(API_ENDPOINTS.AUTH.TOKEN);
        const responseData = res.data;

        if (!responseData.data?.valid) {
          console.error("Token validation failed:", responseData.message);
          message.error("Login failed. Invalid token provided.");
          dispatch(logout()); // Clear invalid token
          navigate("/login");
        } else {
          // Token is valid, complete the login
          message.success("Login successful!");
          dispatch(
            loginSuccess({
              user: responseData.data.user,
              token,
              remember: false,
            })
          );
          navigate("/");
        }
      } catch (error) {
        console.error("Token validation error:", error);
        await handleError(error as Error);
        dispatch(logout()); // Clear token on error
        navigate("/login");
      }
    };

    handleLogin();
  }, [dispatch, searchParams, navigate, handleError]);

  return <LoadingOverlay loading={true} />;
};

export default CallBack;
