package com.seamlance.perfume.controller;

import com.seamlance.perfume.delegate.DiscountDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.DiscountInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/discount")
public class DiscountController {
    @Autowired
    private DiscountDelegate discountDelegate;


    @PostMapping()
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<List<DiscountInfo>> addDiscount(@RequestBody List<DiscountInfo> discountInfos) {
        return ApiResponse.success(discountDelegate.addDiscount(discountInfos));
    }

    @PutMapping
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<List<DiscountInfo>> updateDiscount(@RequestBody Map<String, Integer> discounts) {
        return ApiResponse.success(discountDelegate.updateDiscount(discounts));
    }


    @GetMapping()
    public ApiResponse<List<DiscountInfo>> getAllProductDiscount() {
        return ApiResponse.success(discountDelegate.getAllProductDiscount());
    }

    @DeleteMapping("/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> deleteDiscount(@PathVariable String id) {
        discountDelegate.deleteDiscount(id);
        return ApiResponse.success("DELETED");
    }
}
