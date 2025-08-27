package com.seamlance.perfume.controller;


import com.seamlance.perfume.delegate.ProductDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.ProductInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

@RestController
@RequestMapping("/api/v1/product")
public class ProductController {
    @Autowired
    private ProductDelegate productDelegate;


    @PostMapping("/add")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<ProductInfo> createProduct(@RequestBody ProductInfo productInfo) {
        return ApiResponse.success(productDelegate.addProduct(productInfo));
    }


    @PutMapping("/update/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<ProductInfo> updateProduct(@PathVariable String id, @RequestBody ProductInfo productInfo) {
        return ApiResponse.success(productDelegate.updateProduct(id, productInfo));
    }

    @GetMapping("/all")
    public ApiResponse<List<ProductInfo>> getAll() {
        return ApiResponse.success(productDelegate.getAllProduct());
    }

    @GetMapping("/get/{id}")
    public ApiResponse<ProductInfo> getProduct(@PathVariable String id) {
        return ApiResponse.success(productDelegate.getProductById(id));
    }

    @DeleteMapping("/delete/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> deleteProduct(@PathVariable String id) {
        productDelegate.deleteProduct(id);
        return ApiResponse.success("DELETED");
    }

    @GetMapping("/get-by-category/{id}")
    public ApiResponse<List<ProductInfo>> getProductByCategory(@PathVariable String id) {
        return ApiResponse.success(productDelegate.getProductByCategory(id));
    }

    @GetMapping("/search")
    public ApiResponse<List<ProductInfo>> searchProduct(@RequestParam String value) {
        return ApiResponse.success(productDelegate.searchProduct(value));
    }

    @PutMapping("/upload-image/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<ProductInfo> uploadImage(@PathVariable String id, @RequestParam MultipartFile file) {
        return ApiResponse.success(productDelegate.uploadImage(id, file));
    }
}
