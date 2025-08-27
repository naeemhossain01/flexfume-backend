package com.seamlance.perfume.config;

import com.seamlance.perfume.audit.ApplicationAuditAware;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.domain.AuditorAware;
import org.springframework.web.client.RestTemplate;

@Configuration
public class AppConfig {
    @Bean
    public AuditorAware<String> auditorAware() {
        return new ApplicationAuditAware();
    }

    @Bean
    public RestTemplate restTemplate() {
        return new RestTemplate();
    }
}
