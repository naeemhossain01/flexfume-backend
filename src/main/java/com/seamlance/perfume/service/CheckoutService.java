package com.seamlance.perfume.service;

import com.seamlance.perfume.dto.CheckoutOtpResponse;
import com.seamlance.perfume.dto.CheckoutOtpVerifyRequest;

public interface CheckoutService {
    CheckoutOtpResponse verifyOtpAndHandleUser(CheckoutOtpVerifyRequest request);
}
