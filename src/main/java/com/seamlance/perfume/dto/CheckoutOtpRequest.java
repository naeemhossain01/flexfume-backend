package com.seamlance.perfume.dto;

import lombok.Data;

@Data
public class CheckoutOtpRequest {
    private String phoneNumber;
    private String name;
    private String email;
    private String buildingName;
    private String road;
    private String area;
    private String city;
}
