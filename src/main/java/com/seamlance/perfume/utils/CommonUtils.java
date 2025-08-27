package com.seamlance.perfume.utils;

import java.util.concurrent.ThreadLocalRandom;

import org.springframework.stereotype.Component;

import com.seamlance.perfume.constants.Constant;

@Component
public class CommonUtils {
    private static final String ALPHABETS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
    private static final String NUMBERS = "0123456789";


    public String generateRandomNumberAndAlphabet(int len, boolean isOnlyNumber, boolean isOnlyAlphabet) {
        StringBuilder pool = new StringBuilder();

        if(isOnlyAlphabet == isOnlyNumber) {
            pool.append(ALPHABETS);
            pool.append(NUMBERS);
        } else {
            pool.append(isOnlyAlphabet ? ALPHABETS : NUMBERS);
        }

        StringBuilder result = new StringBuilder(len);
        ThreadLocalRandom random = ThreadLocalRandom.current();

        for(int i = 0; i < len; i++) {
            int index = random.nextInt(pool.length());
            result.append(pool.charAt(index));
        }

        return result.toString();
    }

    public String generateOTPSMS(String otp) {
        StringBuilder pool = new StringBuilder();
        pool.append(Constant.OTP_SMS_TITLE);
        pool.append(Constant.OTP_SMS_BODY + otp + '\n');
        pool.append(Constant.OTP_SMS_FOOTER);

        return pool.toString();
    }

    public String generateResetPasswordSms(String otp) {
        StringBuilder pool = new StringBuilder();
        pool.append(Constant.PASSWORD_RESET_SMS_BODY + otp + '\n');
        pool.append(Constant.PASSWORD_RESET_SMS_FOOTER);

        return pool.toString();
    }

    public String generateCheckoutOTPSMS(String otp) {
        StringBuilder pool = new StringBuilder();
        pool.append("FlexFume Checkout Verification\n");
        pool.append("Your OTP for order verification is " + otp + '\n');
        pool.append("OTP will expire in 5 minutes.");

        return pool.toString();
    }
}
