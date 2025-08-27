package com.seamlance.perfume.dto;

import lombok.Data;

@Data
public class CheckoutOtpVerifyRequest {
    private String phoneNumber;
    private String otp;
    private String name;
    private String email;
    private String buildingName;
    private String road;
    private String area;
    private String city;
}
