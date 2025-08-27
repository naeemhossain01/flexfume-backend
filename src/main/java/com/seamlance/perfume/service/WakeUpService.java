package com.seamlance.perfume.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.HashMap;
import java.util.Map;

@Service
public class WakeUpService {

    @Autowired
    private DataSource dataSource;

    @Autowired
    private RedisTemplate<String, Object> redisTemplate;

    public Map<String, Object> performWakeUpChecks() {
        Map<String, Object> response = new HashMap<>();
        Map<String, Object> checks = new HashMap<>();
        
        // Current timestamp
        String timestamp = LocalDateTime.now().format(DateTimeFormatter.ISO_LOCAL_DATE_TIME);
        response.put("timestamp", timestamp);
        response.put("message", "Wake-up checks completed");
        
        // Check database connection
        checks.put("database", checkDatabase());
        
        // Check Redis connection
        checks.put("redis", checkRedis());
        
        // Check application health
        checks.put("application", checkApplication());
        
        response.put("checks", checks);
        
        // Overall status
        boolean allHealthy = checks.values().stream()
            .allMatch(check -> check instanceof Map && "UP".equals(((Map<?, ?>) check).get("status")));
        
        response.put("status", allHealthy ? "UP" : "PARTIAL");
        response.put("wakeUpTime", System.currentTimeMillis());
        
        return response;
    }

    private Map<String, Object> checkDatabase() {
        Map<String, Object> result = new HashMap<>();
        long startTime = System.currentTimeMillis();
        
        try (Connection connection = dataSource.getConnection()) {
            // Simple query to wake up the database
            try (PreparedStatement stmt = connection.prepareStatement("SELECT 1 as test_connection")) {
                try (ResultSet rs = stmt.executeQuery()) {
                    if (rs.next() && rs.getInt("test_connection") == 1) {
                        result.put("status", "UP");
                        result.put("message", "Database connection successful");
                    } else {
                        result.put("status", "DOWN");
                        result.put("message", "Database query failed");
                    }
                }
            }
        } catch (Exception e) {
            result.put("status", "DOWN");
            result.put("message", "Database connection failed: " + e.getMessage());
        }
        
        long responseTime = System.currentTimeMillis() - startTime;
        result.put("responseTimeMs", responseTime);
        
        return result;
    }

    private Map<String, Object> checkRedis() {
        Map<String, Object> result = new HashMap<>();
        long startTime = System.currentTimeMillis();
        
        try {
            // Test Redis connection with a simple ping
            String testKey = "wake_up_test";
            String testValue = "ping_" + System.currentTimeMillis();
            
            redisTemplate.opsForValue().set(testKey, testValue);
            String retrievedValue = (String) redisTemplate.opsForValue().get(testKey);
            
            if (testValue.equals(retrievedValue)) {
                result.put("status", "UP");
                result.put("message", "Redis connection successful");
                
                // Clean up test key
                redisTemplate.delete(testKey);
            } else {
                result.put("status", "DOWN");
                result.put("message", "Redis value mismatch");
            }
        } catch (Exception e) {
            result.put("status", "DOWN");
            result.put("message", "Redis connection failed: " + e.getMessage());
        }
        
        long responseTime = System.currentTimeMillis() - startTime;
        result.put("responseTimeMs", responseTime);
        
        return result;
    }

    private Map<String, Object> checkApplication() {
        Map<String, Object> result = new HashMap<>();
        
        // Check JVM memory
        Runtime runtime = Runtime.getRuntime();
        long maxMemory = runtime.maxMemory();
        long totalMemory = runtime.totalMemory();
        long freeMemory = runtime.freeMemory();
        long usedMemory = totalMemory - freeMemory;
        
        Map<String, Object> memory = new HashMap<>();
        memory.put("maxMemoryMB", maxMemory / (1024 * 1024));
        memory.put("totalMemoryMB", totalMemory / (1024 * 1024));
        memory.put("usedMemoryMB", usedMemory / (1024 * 1024));
        memory.put("freeMemoryMB", freeMemory / (1024 * 1024));
        memory.put("memoryUsagePercent", (usedMemory * 100.0) / totalMemory);
        
        result.put("status", "UP");
        result.put("message", "Application is running");
        result.put("memory", memory);
        result.put("availableProcessors", runtime.availableProcessors());
        
        return result;
    }
}
