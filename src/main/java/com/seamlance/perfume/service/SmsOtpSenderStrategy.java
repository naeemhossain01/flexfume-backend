package com.seamlance.perfume.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.seamlance.perfume.constants.Constant;
import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.enums.OtpSenderType;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.utils.CommonUtils;
import com.seamlance.perfume.utils.ValidationUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

import java.util.HashMap;
import java.util.Map;

@Service
public class SmsOtpSenderStrategy implements OtpSenderStrategy {

    private final RestTemplate restTemplate = new RestTemplate();

    @Value("${otp.sms.url}")
    private String url;

    @Value("${otp.sms.apiKey}")
    private String apiKey;

    @Value("${otp.sms.senderId}")
    private String senderId;

    @Autowired
    private CommonUtils commonUtils;

    @Autowired
    private RedisService redisService;

    @Override
    public void sendOtp(String number, String textMessage) {

        if(isOtpStillValid(number)) {
            throw new InvalidRequestsException(Constant.OTP_ALREADY_SEND);
        }


        JsonNode response = null;
        try {
            response = sendSms(number, textMessage);
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        ValidationUtils.otpResponseValidation(response);
    }

    @Override
    public OtpSenderType getType() {
        return OtpSenderType.SMS;
    }

    private String getSmsPayload(String number, String message) throws JsonProcessingException {
        ObjectMapper mapper = new ObjectMapper();


        Map<String, String> payload = new HashMap<>();
        payload.put(Constant.OTP_SMS_API_KEY, apiKey);
        payload.put(Constant.OTP_SMS_SENDER_ID, senderId);
        payload.put(Constant.OTP_SMS_CLIENT_NUMBER_KEY, number);
        payload.put(Constant.OTP_SMS_CLIENT_MESSAGE_KEY, message);

        String jsonPayload = mapper.writeValueAsString(payload);

        return jsonPayload;
    }

    private boolean isOtpStillValid(String number) {
        String existingOtp = redisService.get(Constant.REDIS_OTP_PREFIX_KEY + number, String.class);

        return existingOtp != null;
    }

    private JsonNode sendSms(String number, String textMessage) throws JsonProcessingException {
        JsonNode response = null;
        String payload = getSmsPayload(number, textMessage);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);

        try {
            HttpEntity<String> request = new HttpEntity<>(payload, headers);

            response = restTemplate.postForObject(url, request, JsonNode.class);
        } catch (Exception e) {
            throw new InvalidRequestsException(ErrorConstant.FAILED_TO_SEND_SMS);
        }

        return response;
    }
}
