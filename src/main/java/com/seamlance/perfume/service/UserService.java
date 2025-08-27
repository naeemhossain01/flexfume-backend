package com.seamlance.perfume.service;

import java.util.List;

import com.seamlance.perfume.dto.ChangePasswordRequest;
import com.seamlance.perfume.dto.LoginRequest;
import com.seamlance.perfume.entity.User;

public interface UserService {
    User registerUser(User user);
    User getUser(String userId);
    List<User> getAllUser();
    User updateUser(String userId, User user);
    String loginUser(LoginRequest loginRequest);
    User getLoginUser();
    void validateAlreadyHaveAccount(String phoneNumber);
    void changePassword(ChangePasswordRequest changePasswordRequest);
    User getUserByPhoneNumber(String phoneNumber);
    void resetPassword(String phoneNumber, String newPassword, String confirmPassword);
}
