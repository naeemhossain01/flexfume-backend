package com.seamlance.perfume.service;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dto.OtpRequest;
import com.seamlance.perfume.enums.OtpSenderType;
import com.seamlance.perfume.exception.InvalidOtpSenderTypeException;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import com.seamlance.perfume.utils.CommonUtils;

@Service
public class OtpService {
    private Map<OtpSenderType, OtpSenderStrategy> strategies = new HashMap<>();

    @Autowired
    private RedisService redisService;

    @Autowired
    private CommonUtils commonUtils;

    public OtpService(List<OtpSenderStrategy> senderStrategyList) {
        strategies = senderStrategyList.stream().collect(Collectors.toMap(OtpSenderStrategy::getType, strategy -> strategy));
    }

    public void send(String phoneNumber, String type, String smsType) {
        OtpSenderStrategy otpSenderStrategy = strategies.get(OtpSenderType.valueOf(type.toUpperCase()));

        if(otpSenderStrategy == null) {
            throw new InvalidOtpSenderTypeException(ErrorConstant.UNSUPPORTED_OTP_SENDER_TYPE + type);
        }

        String otp = commonUtils.generateRandomNumberAndAlphabet(4, true, false);
        String textMessage = getSms(smsType, otp);

        otpSenderStrategy.sendOtp(phoneNumber, textMessage);

        storeOtpInCache(phoneNumber, otp, smsType);
    }

    public void verifyOtpForUser(OtpRequest otpRequest) {
        if(otpRequest == null) {
            throw new InvalidRequestsException(ErrorConstant.INTERNAL_SERVER_ERROR);
        }

        String cachedOtp = null;

        try {
            cachedOtp = redisService.get(Constant.REDIS_OTP_PREFIX_KEY + otpRequest.getPhoneNumber(), String.class);

            if(cachedOtp == null || !cachedOtp.equals(otpRequest.getOtp())) {
                throw new InvalidRequestsException(ErrorConstant.INVALID_OTP);
            }
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    public void markPhoneNumberAsValid(OtpRequest otpRequest) {
        otpRequest.setVerified(true);
        redisService.set(Constant.REDIS_PHONE_VERIFIED_KEY + otpRequest.getPhoneNumber(), otpRequest, 1800L);
    }

    public void isNumberIsRegistered(String number) {

        if(number == null) {
            throw new ResourceNotFoundException(ErrorConstant.PHONE_NUMBER_NOT_FOUND);
        }

        OtpRequest cachedOtpRequest = null;

        try {
            cachedOtpRequest = redisService.get(Constant.REDIS_PHONE_VERIFIED_KEY + number, OtpRequest.class);

            if(cachedOtpRequest == null) {
                throw new ResourceNotFoundException(ErrorConstant.TIME_EXPIRED);
            }

            if(!cachedOtpRequest.isVerified()) {
                throw new InvalidRequestsException(ErrorConstant.PHONE_NUMBER_NOT_REGISTERED);
            }

        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    public void isPhoneNumberVerified(String phoneNumber) {
        if(phoneNumber == null) {
            throw new InvalidRequestsException(ErrorConstant.INVALID_NUMBER);
        }

        try {
            OtpRequest otpRequest = redisService.get(Constant.REDIS_PHONE_VERIFIED_KEY + phoneNumber, OtpRequest.class);
            if(otpRequest != null && otpRequest.isVerified()) {
                throw new InvalidRequestsException(ErrorConstant.PHONE_NUMBER_ALREADY_VERIFIED);
            }
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    public void verifyResetPasswordOtp(String phoneNumber, String otp) {
        String cachedOtp = redisService.get(Constant.REDIS_OTP_RESET_PASSWORD_KEY + phoneNumber, String.class);

        if(cachedOtp == null || !cachedOtp.equals(otp)) {
            throw new InvalidRequestsException(ErrorConstant.INVALID_OTP);
        }
    }

    private String getSms(String smsType, String otp) {
        String sms = null;

        switch (smsType) {
            case Constant.USER_REGISTRATION_TYPE -> {
                sms = commonUtils.generateOTPSMS(otp);
                break;
            }
            case Constant.USER_PASSWORD_RESET_TYPE -> {
                sms = commonUtils.generateResetPasswordSms(otp);
                break;
            }
            case Constant.USER_CHECKOUT_TYPE -> {
                sms = commonUtils.generateCheckoutOTPSMS(otp);
                break;
            }
            default -> {
               throw new InvalidRequestsException(ErrorConstant.INVALID_SMS_TYPE);
            }
        }

        return sms;
    }

    private void storeOtpInCache(String number, String otp, String smsType) {
        String prefixKey = Constant.REDIS_OTP_PREFIX_KEY;
        long expiredTime = 300L;

        if(smsType.equalsIgnoreCase(Constant.USER_PASSWORD_RESET_TYPE)) {
            prefixKey = Constant.REDIS_OTP_RESET_PASSWORD_KEY;
            expiredTime = 600L;
        }

        redisService.set(prefixKey + number, otp, expiredTime);
    }
}
