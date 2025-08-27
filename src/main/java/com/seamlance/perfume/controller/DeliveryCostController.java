package com.seamlance.perfume.controller;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import com.seamlance.perfume.delegate.DeliveryCostDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.DeliveryCostInfo;

@RestController
@RequestMapping("/api/v1/delivery-cost")
public class DeliveryCostController {
    @Autowired
    private DeliveryCostDelegate deliveryCostDelegate;

    @PostMapping()
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<DeliveryCostInfo> addCost(@RequestBody DeliveryCostInfo deliveryCostInfo) {
        return ApiResponse.success(deliveryCostDelegate.addCost(deliveryCostInfo));
    }

    @PutMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<DeliveryCostInfo> updateCost(@PathVariable String id, @RequestBody DeliveryCostInfo deliveryCostInfo) {
        return ApiResponse.success(deliveryCostDelegate.updateCost(id, deliveryCostInfo));
    }

    @GetMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<DeliveryCostInfo> getDeliveryCost(@PathVariable String id) {
        return ApiResponse.success(deliveryCostDelegate.getDeliveryCostById(id));
    }

    @GetMapping("/all")
    public ApiResponse<List<DeliveryCostInfo>> getAllDeliveryCost() {
        return ApiResponse.success(deliveryCostDelegate.getAllDeliveryCost());
    }

    @GetMapping("/location")
    public ApiResponse<List<DeliveryCostInfo>> getDeliveryCostByLocation(@RequestParam String location) {
        return ApiResponse.success(deliveryCostDelegate.getDeliveryCostByLocation(location));
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> deleteDeliveryCost(@PathVariable String id) {
        deliveryCostDelegate.deleteDeliveryCost(id);
        return ApiResponse.success("Delivery cost deleted");
    }
}
