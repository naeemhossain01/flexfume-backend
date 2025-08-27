package com.seamlance.perfume.service;

import com.seamlance.perfume.enums.OtpSenderType;

public interface OtpSenderStrategy {
    void sendOtp(String number, String textMessage);
    OtpSenderType getType();
}
