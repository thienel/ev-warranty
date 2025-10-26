import React from "react";
import { Button, Typography, Progress } from "antd";
import { ArrowLeftOutlined } from "@ant-design/icons";

const { Title, Text } = Typography;

interface ClaimFormHeaderProps {
  onBack: () => void;
  currentStep: number;
  totalSteps: number;
}

const ClaimFormHeader: React.FC<ClaimFormHeaderProps> = ({
  onBack,
  currentStep,
  totalSteps,
}) => {
  const progressPercent = Math.round(((currentStep + 1) / totalSteps) * 100);

  return (
    <div className="claim-create-header">
      <div className="header-nav">
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={onBack}
          type="text"
          size="large"
          className="back-button"
        >
          Back to Claims
        </Button>
      </div>

      <div className="header-content">
        <div className="header-text">
          <Title level={1} className="page-title">
            Create New Warranty Claim
          </Title>
          <Text className="page-description">
            Create a new warranty claim for a customer's vehicle. Follow the
            steps below to complete the process.
          </Text>
        </div>

        <div className="header-progress">
          <Progress
            type="circle"
            percent={progressPercent}
            size={80}
            strokeColor={{
              "0%": "#697565",
              "100%": "#5a6358",
            }}
            format={() => (
              <div className="progress-content">
                <div className="step-number">{currentStep + 1}</div>
                <div className="step-total">of {totalSteps}</div>
              </div>
            )}
          />
        </div>
      </div>
    </div>
  );
};

export default ClaimFormHeader;
