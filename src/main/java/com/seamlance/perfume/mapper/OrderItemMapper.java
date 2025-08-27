package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.OrderItem;
import com.seamlance.perfume.info.OrderInfo;
import com.seamlance.perfume.info.OrderItemInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {ProductMapper.class})
public interface OrderItemMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "productInfo", target = "product")
    OrderItem toEntity(OrderItemInfo orderItemInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "productInfo", source = "product")
    OrderItemInfo toInfo(OrderItem orderItem);

    List<OrderItemInfo> toInfoList(List<OrderItem> orderItemList);
}
