package com.seamlance.perfume.audit;

import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.security.AuthUser;
import org.springframework.data.domain.AuditorAware;
import org.springframework.security.authentication.AnonymousAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

import java.util.Optional;

public class ApplicationAuditAware implements AuditorAware<String> {
    @Override
    public Optional<String> getCurrentAuditor() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();

        if(authentication == null || !authentication.isAuthenticated() || authentication instanceof AnonymousAuthenticationToken) {
            return Optional.empty();
        }

        AuthUser userPrinciple = (AuthUser) authentication.getPrincipal();
        return Optional.ofNullable(userPrinciple.getUsername());
    }
}
