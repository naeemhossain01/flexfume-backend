package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.DeliveryCost;
import com.seamlance.perfume.info.DeliveryCostInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface DeliveryCostMapper {

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    DeliveryCost toEntity(DeliveryCostInfo deliveryCostInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    DeliveryCostInfo toInfo(DeliveryCost deliveryCost);

    List<DeliveryCostInfo> toInfoList(List<DeliveryCost> deliveryCostList);
}
