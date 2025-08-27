package com.seamlance.perfume.service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.seamlance.perfume.exception.InvalidRequestsException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.util.concurrent.TimeUnit;

@Service
public class RedisService {

    @Autowired
    private RedisTemplate redisTemplate;

    public <T> T get(String key, Class<T> entityClass) {
        ObjectMapper mapper = new ObjectMapper();

        try {
            Object o = redisTemplate.opsForValue().get(key);

            if(o == null) {
                return null;
            }

            return mapper.readValue(o.toString(), entityClass);
        } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }

    public void set(String key, Object o, Long ttl) {
        ObjectMapper objectMapper = new ObjectMapper();

        try {
            String jsonValue = objectMapper.writeValueAsString(o);
            redisTemplate.opsForValue().set(key, jsonValue, ttl, TimeUnit.SECONDS);
         } catch (Exception e) {
            throw new InvalidRequestsException(e.getMessage());
        }
    }
}
