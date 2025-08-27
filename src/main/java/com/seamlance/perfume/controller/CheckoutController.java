package com.seamlance.perfume.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.dto.CheckoutOtpRequest;
import com.seamlance.perfume.dto.CheckoutOtpResponse;
import com.seamlance.perfume.dto.CheckoutOtpVerifyRequest;
import com.seamlance.perfume.service.CheckoutService;
import com.seamlance.perfume.service.OtpService;

@RestController
@RequestMapping("/api/v1/checkout")
public class CheckoutController {

    @Autowired
    private OtpService otpService;

    @Autowired
    private CheckoutService checkoutService;

    @PostMapping("/send-otp")
    public ApiResponse<String> sendCheckoutOtp(@RequestBody CheckoutOtpRequest request) {
        // Send OTP for checkout verification
        otpService.send(request.getPhoneNumber(), "SMS", Constant.USER_CHECKOUT_TYPE);
        return ApiResponse.success("OTP sent successfully for checkout verification");
    }

    @PostMapping("/verify-otp")
    public ApiResponse<CheckoutOtpResponse> verifyCheckoutOtp(@RequestBody CheckoutOtpVerifyRequest request) {
        try {
            // Verify OTP and handle user account creation/update
            CheckoutOtpResponse response = checkoutService.verifyOtpAndHandleUser(request);
            return ApiResponse.success(response);
        } catch (Exception e) {
            // Handle any exceptions and return proper error response
            return ApiResponse.error("Checkout verification failed: " + e.getMessage());
        }
    }
}
