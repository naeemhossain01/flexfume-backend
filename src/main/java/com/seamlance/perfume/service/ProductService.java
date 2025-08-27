package com.seamlance.perfume.service;

import com.seamlance.perfume.entity.Product;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

public interface ProductService {
    Product addProduct(Product product);
    Product updateProduct(String productId, Product product);
    Product getProductById(String productId);
    List<Product> getAllProduct();
    void deleteProduct(String productId);
    List<Product> getProductByCategory(String categoryId);
    List<Product> searchProduct(String value);
    Product uploadImage(String productId, MultipartFile file);
}
