package com.seamlance.perfume.delegate;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import com.seamlance.perfume.dto.ChangePasswordRequest;
import com.seamlance.perfume.dto.LoginRequest;
import com.seamlance.perfume.info.UserInfo;
import com.seamlance.perfume.mapper.UserMapper;
import com.seamlance.perfume.service.UserService;

@Component
public class UserDelegateImpl implements UserDelegate{
    @Autowired
    private UserService userService;

    @Autowired
    private UserMapper userMapper;


    @Override
    public UserInfo registerUser(UserInfo userInfo) {
        return userMapper.toInfo(userService.registerUser(userMapper.toEntity(userInfo)));
    }

    @Override
    public UserInfo getUser(String userId) {
        return userMapper.toInfo(userService.getUser(userId));
    }

    @Override
    public UserInfo getCurrentUserProfile() {
        return userMapper.toInfo(userService.getLoginUser());
    }

    @Override
    public List<UserInfo> getAllUser() {
        return userMapper.toInfoList(userService.getAllUser());
    }

    @Override
    public UserInfo updateUser(String userId, UserInfo userInfo) {
        return userMapper.toInfo(userService.updateUser(userId, userMapper.toEntity(userInfo)));
    }

    @Override
    public String loginUser(LoginRequest loginRequest) {
        return userService.loginUser(loginRequest);
    }

    @Override
    public void validateAlreadyHaveAccount(String phoneNumber) {
        userService.validateAlreadyHaveAccount(phoneNumber);
    }

    @Override
    public void changePassword(ChangePasswordRequest changePasswordRequest) {
        userService.changePassword(changePasswordRequest);
    }

    @Override
    public UserInfo getUserByPhoneNumber(String phoneNumber) {
        return userMapper.toInfo(userService.getUserByPhoneNumber(phoneNumber));
    }

    @Override
    public void resetPassword(String phoneNumber, String newPassword, String confirmPassword) {
        userService.resetPassword(phoneNumber, newPassword, confirmPassword);
    }
}
