package com.seamlance.perfume.delegate;

import java.util.List;

import com.seamlance.perfume.dto.ChangePasswordRequest;
import com.seamlance.perfume.dto.LoginRequest;
import com.seamlance.perfume.info.UserInfo;

public interface UserDelegate {
    UserInfo registerUser(UserInfo userInfo);
    UserInfo getUser(String userId);
    UserInfo getCurrentUserProfile();
    List<UserInfo> getAllUser();
    UserInfo updateUser(String userId, UserInfo userInfo);
    String loginUser(LoginRequest loginRequest);
    void validateAlreadyHaveAccount(String phoneNumber);
    void changePassword(ChangePasswordRequest changePasswordRequest);
    UserInfo getUserByPhoneNumber(String phoneNumber);
    void resetPassword(String phoneNumber, String newPassword, String confirmPassword);
}
