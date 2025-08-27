package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.Cart;
import com.seamlance.perfume.info.CartInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {ProductMapper.class, UserMapper.class})
public interface CartMapper {

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "productInfo", target = "product")
    @Mapping(source = "userInfo", target = "user")
    Cart toEntity(CartInfo cartInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "productInfo", source = "product")
    @Mapping(target = "userInfo", source = "user")
    CartInfo toInfo(Cart cart);

    List<CartInfo> toInfoList(List<Cart> carts);
    List<Cart> toEntityList(List<CartInfo> cartInfoList);
}
