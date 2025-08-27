package com.seamlance.perfume.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import com.seamlance.perfume.dao.UserDao;
import com.seamlance.perfume.dto.CheckoutOtpResponse;
import com.seamlance.perfume.dto.CheckoutOtpVerifyRequest;
import com.seamlance.perfume.dto.OtpRequest;
import com.seamlance.perfume.entity.Address;
import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.security.JwtUtils;
import com.seamlance.perfume.utils.CommonUtils;

@Service
public class CheckoutServiceImpl implements CheckoutService {

    @Autowired
    private OtpService otpService;

    @Autowired
    private UserDao userDao;



    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private JwtUtils jwtUtils;

    @Autowired
    private CommonUtils commonUtils;

    @Override
    @Transactional
    public CheckoutOtpResponse verifyOtpAndHandleUser(CheckoutOtpVerifyRequest request) {
        // First verify the OTP
        OtpRequest otpRequest = new OtpRequest();
        otpRequest.setPhoneNumber(request.getPhoneNumber());
        otpRequest.setOtp(request.getOtp());
        
        otpService.verifyOtpForUser(otpRequest);

        // Check if user exists with this phone number
        User existingUser = userDao.findUserByPhoneNumber(request.getPhoneNumber()).orElse(null);
        
        CheckoutOtpResponse response = new CheckoutOtpResponse();
        
        if (existingUser != null) {
            // User exists - update their information
            existingUser.setName(request.getName());
            existingUser.setEmail(request.getEmail());
            
            // Update address
            Address address = existingUser.getAddress();
            if (address == null) {
                address = new Address();
                address.setUser(existingUser);
                // Set audit fields for new address
                address.setCreatedBy("SYSTEM");
                address.setCreatedDate(java.time.LocalDateTime.now());
            }
            address.setBuildingName(request.getBuildingName());
            address.setRoad(request.getRoad());
            address.setArea(request.getArea());
            address.setCity(request.getCity());
            
            existingUser.setAddress(address);
            existingUser = userDao.save(existingUser);
            
            response.setNewUser(false);
            response.setUserId(existingUser.getId());
            response.setUserName(existingUser.getName());
            response.setUserEmail(existingUser.getEmail());
            response.setUserPhoneNumber(existingUser.getPhoneNumber());
            
            // Generate JWT token
            String token = jwtUtils.generateToken(existingUser.getPhoneNumber(), existingUser.getRole());
            response.setToken(token);
            
        } else {
            // Create new user
            User newUser = new User();
            newUser.setName(request.getName());
            newUser.setEmail(request.getEmail());
            newUser.setPhoneNumber(request.getPhoneNumber());
            
            // Generate a random password for the user
            String tempPassword = commonUtils.generateRandomNumberAndAlphabet(12, true, true);
            newUser.setPassword(passwordEncoder.encode(tempPassword));
            newUser.setRole("USER");
            
            // Set audit fields manually for system-created user during checkout
            newUser.setCreatedBy("SYSTEM");
            newUser.setCreatedDate(java.time.LocalDateTime.now());
            
            // Create address
            Address address = new Address();
            address.setBuildingName(request.getBuildingName());
            address.setRoad(request.getRoad());
            address.setArea(request.getArea());
            address.setCity(request.getCity());
            address.setUser(newUser);
            
            // Set audit fields manually for system-created address during checkout
            address.setCreatedBy("SYSTEM");
            address.setCreatedDate(java.time.LocalDateTime.now());
            
            newUser.setAddress(address);
            newUser = userDao.save(newUser);
            
            response.setNewUser(true);
            response.setUserId(newUser.getId());
            response.setUserName(newUser.getName());
            response.setUserEmail(newUser.getEmail());
            response.setUserPhoneNumber(newUser.getPhoneNumber());
            
            // Generate JWT token
            String token = jwtUtils.generateToken(newUser.getPhoneNumber(), newUser.getRole());
            response.setToken(token);
        }
        
        return response;
    }
}
