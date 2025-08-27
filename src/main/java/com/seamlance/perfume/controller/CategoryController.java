package com.seamlance.perfume.controller;

import com.seamlance.perfume.delegate.CategoryDelegate;
import com.seamlance.perfume.dto.ApiResponse;
import com.seamlance.perfume.info.CategoryInfo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/category")
public class CategoryController {

    @Autowired
    private CategoryDelegate categoryDelegate;

    @PostMapping("/add")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<CategoryInfo> createCategory(@RequestBody CategoryInfo categoryInfo) {
        return ApiResponse.success(categoryDelegate.createCategory(categoryInfo));
    }

    @PutMapping("/update/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<CategoryInfo> updateCategory(@PathVariable String id, @RequestBody CategoryInfo categoryInfo) {
        return ApiResponse.success(categoryDelegate.updateCategory(id, categoryInfo));
    }

    @GetMapping("/all")
    public ApiResponse<List<CategoryInfo>> getAll() {
        return ApiResponse.success(categoryDelegate.getAllCategories());
    }

    @GetMapping("/get/{id}")
    public ApiResponse<CategoryInfo> getCategory(@PathVariable String id) {
        return ApiResponse.success(categoryDelegate.getCategoryById(id));
    }

    @DeleteMapping("/delete/{id}")
    @PreAuthorize("hasAuthority('ADMIN')")
    public ApiResponse<String> deleteCategory(@PathVariable String id) {
        categoryDelegate.deleteCategory(id);
        return ApiResponse.success("Deleted");
    }
}
