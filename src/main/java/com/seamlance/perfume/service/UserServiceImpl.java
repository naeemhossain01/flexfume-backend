package com.seamlance.perfume.service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DataAccessException;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.web.client.ResourceAccessException;

import com.seamlance.perfume.constants.EntityConstant;
import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.UserDao;
import com.seamlance.perfume.dto.ChangePasswordRequest;
import com.seamlance.perfume.dto.LoginRequest;
import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.exception.InvalidCredentialsException;
import com.seamlance.perfume.exception.InvalidRequestsException;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import com.seamlance.perfume.security.JwtUtils;

@Service
public class UserServiceImpl implements UserService {

    @Autowired
    private UserDao userDao;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private JwtUtils jwtUtils;

    @Override
    public User registerUser(User user) {
        //TODO: Need to add validation

        User existingUser = userDao.findUserByPhoneNumber(user.getPhoneNumber()).orElse(null);

        if(existingUser != null) {
            throw new InvalidRequestsException(ErrorConstant.ALREADY_HAVE_ACCOUNT);
        }

        // Set role logic: if ADMIN role is explicitly set, keep it; otherwise default to USER
        if(user.getRole() != null && user.getRole().equalsIgnoreCase(EntityConstant.ADMIN_ROLE)) {
            user.setRole(EntityConstant.ADMIN_ROLE);
        } else {
            // Set default role to USER for all new registrations
            user.setRole(EntityConstant.USER_ROLE);
        }

        try {
            user.setPassword(passwordEncoder.encode(user.getPassword()));
            user.setCreatedBy(user.getPhoneNumber()); // Use phone number instead of email
            user = userDao.saveAndFlush(user);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }

        return user;
    }

    @Override
    public User getUser(String userId) {
        return userDao.findById(userId).orElseThrow(() -> new ResourceAccessException(ErrorConstant.USER_NOT_FOUND));
    }

    @Override
    public List<User> getAllUser() {
        List<User> users = new ArrayList<>();
        try {
            users = userDao.findAll();
        } catch (DataAccessException e) {
            e.printStackTrace();;
        }

        return users;
    }

    @Override
    public User updateUser(String userId, User user) {
        //TODO: need validation here

        User updatedUser = userDao.findById(userId).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.USER_NOT_FOUND));

        if(user.getName() != null) updatedUser.setName(user.getName());
        if(user.getEmail() != null) updatedUser.setEmail(user.getEmail());
        if(user.getPhoneNumber() != null) updatedUser.setPhoneNumber(user.getPhoneNumber());
        if(user.getAddress() != null) updatedUser.setAddress(user.getAddress());

        try {
            updatedUser = userDao.save(updatedUser);
        } catch (DataAccessException e) {
            e.printStackTrace();
        }

        return updatedUser;
    }

    @Override
    public String loginUser(LoginRequest loginRequest) {
        User user = userDao.findUserByPhoneNumber(loginRequest.getPhoneNumber()).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.USER_NOT_FOUND));

        if(!passwordEncoder.matches(loginRequest.getPassword(), user.getPassword())) {
            throw new InvalidCredentialsException(ErrorConstant.BAD_CREDENTIALS);
        }

        String token = jwtUtils.generateToken(user);

        if(token == null) {
            throw new ResourceNotFoundException(ErrorConstant.CANT_GENERATE_TOKEN);
        }

        return token;
    }

    @Override
    public User getLoginUser() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        String userName = authentication.getName();

        return userDao.findUserByPhoneNumber(userName).orElseThrow(() -> new UsernameNotFoundException(ErrorConstant.USER_NOT_FOUND));
    }

    @Override
    public void validateAlreadyHaveAccount(String phoneNumber) {
        User user = userDao.findUserByPhoneNumber(phoneNumber).orElse(null);

        if(user != null) {
            throw new InvalidRequestsException(ErrorConstant.ALREADY_HAVE_ACCOUNT);
        }
    }

    @Override
    public void changePassword(ChangePasswordRequest changePasswordRequest) {
        // Get the currently logged-in user
        User currentUser = getLoginUser();

        // Validate that new password and confirm password match
        if (!changePasswordRequest.getNewPassword().equals(changePasswordRequest.getConfirmPassword())) {
            throw new InvalidRequestsException(ErrorConstant.PASSWORD_MISMATCH);
        }

        // Validate that current password is correct
        if (!passwordEncoder.matches(changePasswordRequest.getCurrentPassword(), currentUser.getPassword())) {
            throw new InvalidCredentialsException(ErrorConstant.CURRENT_PASSWORD_INCORRECT);
        }

        // Validate that new password is different from current password
        if (passwordEncoder.matches(changePasswordRequest.getNewPassword(), currentUser.getPassword())) {
            throw new InvalidRequestsException(ErrorConstant.SAME_PASSWORD_ERROR);
        }

        // Update the password
        currentUser.setPassword(passwordEncoder.encode(changePasswordRequest.getNewPassword()));
        currentUser.setLastModifiedBy(currentUser.getPhoneNumber());

        try {
            userDao.save(currentUser);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    @Override
    public User getUserByPhoneNumber(String phoneNumber) {
        if(phoneNumber == null) {
            throw new InvalidRequestsException(ErrorConstant.PHONE_NUMBER_NOT_FOUND);
        }

        User user = userDao.findUserByPhoneNumber(phoneNumber).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.USER_NOT_FOUND));

        return user;
    }

    @Override
    public void resetPassword(String phoneNumber, String newPassword, String confirmPassword) {
        User user = this.getUserByPhoneNumber(phoneNumber);

        if(!newPassword.equals(confirmPassword)) {
            throw new InvalidRequestsException(ErrorConstant.PASSWORD_MISMATCH);
        }

        try {
            user.setPassword(passwordEncoder.encode(newPassword));
            user.setLastModifiedDate(LocalDateTime.now());
            user.setLastModifiedBy(user.getPhoneNumber());

            userDao.saveAndFlush(user);
        } catch (DataAccessException e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }
}
