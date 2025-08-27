package com.seamlance.perfume.delegate;

import com.seamlance.perfume.info.ProductInfo;
import org.springframework.web.multipart.MultipartFile;

import java.util.List;

public interface ProductDelegate {
    ProductInfo addProduct(ProductInfo productInfo);
    ProductInfo getProductById(String id);
    ProductInfo updateProduct(String id, ProductInfo productInfo);
    List<ProductInfo> getAllProduct();
    void deleteProduct(String id);
    List<ProductInfo> getProductByCategory(String categoryId);
    List<ProductInfo> searchProduct(String value);
    ProductInfo uploadImage(String productId, MultipartFile file);
}
