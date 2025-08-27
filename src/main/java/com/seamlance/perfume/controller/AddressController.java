package com.seamlance.perfume.controller;

import com.seamlance.perfume.delegate.AddressDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.AddressInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/address")
public class AddressController {
    @Autowired
    private AddressDelegate addressDelegate;

    @PostMapping()
    public ApiResponse<AddressInfo> addAddress(@RequestBody AddressInfo addressInfo) {
        return ApiResponse.success(addressDelegate.addAddress(addressInfo));
    }

    @PutMapping("/{id}")
    public ApiResponse<AddressInfo> updateAddress(@PathVariable String id, @RequestBody AddressInfo addressInfo) {
        return  ApiResponse.success(addressDelegate.updateAddress(id, addressInfo));
    }

    @GetMapping("/{userId}")
    public ApiResponse<AddressInfo> getAddress(@PathVariable String userId) {
        return ApiResponse.success(addressDelegate.getAddressByUser(userId));
    }
}
