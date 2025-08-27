package com.seamlance.perfume.dto;

import lombok.Data;

@Data
public class CheckoutOtpResponse {
    private String token;
    private boolean newUser;
    private String userId;
    private String userName;
    private String userEmail;
    private String userPhoneNumber;
}
