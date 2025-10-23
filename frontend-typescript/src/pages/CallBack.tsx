import React, { useEffect } from "react";
import LoadingOverlay from "@components/LoadingOverlay/LoadingOverlay.tsx";
import { useNavigate, useSearchParams } from "react-router-dom";
import { loginSuccess, setToken } from "@redux/authSlice.js";
import { message } from "antd";
import { useDispatch } from "react-redux";
import api from "@services/api.js";
import { API_ENDPOINTS } from "@constants/common-constants.js";
import useHandleApiError from "@/hooks/useHandleApiError.js";

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
        dispatch(setToken(token));
        const res = (
          await api.get(API_ENDPOINTS.AUTH.TOKEN, { withCredentials: true })
        ).data;
        if (!res.data.valid) {
          console.error(res.message);
          message.error("Login failed. Please check your credentials.");
        } else {
          navigate("/");
          dispatch(
            loginSuccess({ user: res.data.user, token, remember: false })
          );
        }
      } catch (error) {
        await handleError(error as Error);
        navigate("/login");
      }
    };

    handleLogin();
  }, [dispatch, searchParams, navigate, handleError]);

  return <LoadingOverlay loading={true} />;
};

export default CallBack;
