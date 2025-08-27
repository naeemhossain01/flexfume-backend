package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.Order;
import com.seamlance.perfume.info.OrderInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = OrderItemMapper.class)
public interface OrderMapper {
    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "orderItemInfoList", target = "orderItemList")
    Order toEntity(OrderInfo orderInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "orderItemInfoList", source = "orderItemList")
    OrderInfo toInfo(Order order);

    List<OrderInfo> toInfoList(List<Order> orders);
}
