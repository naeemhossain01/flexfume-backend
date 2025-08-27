package com.seamlance.perfume.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.seamlance.perfume.service.WakeUpService;

import java.util.Map;

@RestController
@RequestMapping("/api/v1/system")
public class WakeUpController {

    @Autowired
    private WakeUpService wakeUpService;

    @GetMapping("/wake-up")
    public ResponseEntity<Map<String, Object>> wakeUp() {
        Map<String, Object> response = wakeUpService.performWakeUpChecks();
        return ResponseEntity.ok(response);
    }

    @GetMapping("/ping")
    public ResponseEntity<Map<String, Object>> ping() {
        return ResponseEntity.ok(Map.of(
            "status", "pong",
            "timestamp", System.currentTimeMillis(),
            "message", "Backend is awake"
        ));
    }
}
