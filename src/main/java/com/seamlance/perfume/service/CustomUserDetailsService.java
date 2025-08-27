package com.seamlance.perfume.service;

import com.seamlance.perfume.constants.ErrorConstant;
import com.seamlance.perfume.dao.UserDao;
import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.exception.ResourceNotFoundException;
import com.seamlance.perfume.security.AuthUser;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

@Service
public class CustomUserDetailsService implements UserDetailsService {

    @Autowired
    private UserDao userDao;

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        User user = userDao.findUserByPhoneNumber(username).orElseThrow(() -> new ResourceNotFoundException(ErrorConstant.EMAIL_NOT_FOUND));

        return getAuthUser(user);
    }

    private AuthUser getAuthUser(User user) {
        AuthUser authUser = new AuthUser();
        authUser.setUser(user);

        return authUser;
    }
}
