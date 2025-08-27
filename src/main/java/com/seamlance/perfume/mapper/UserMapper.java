package com.seamlance.perfume.mapper;

import com.seamlance.perfume.entity.User;
import com.seamlance.perfume.info.UserInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring", uses = {AddressMapper.class, OrderMapper.class })
public interface UserMapper {

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(source = "createdTime", target = "createdDate")
    @Mapping(source = "updatedBy" , target = "lastModifiedBy")
    @Mapping(source = "updatedTime", target = "lastModifiedDate")
    @Mapping(source = "addressInfo", target = "address")
    @Mapping(source = "orderInfoList", target = "orders")
    User toEntity(UserInfo userInfo);

    @Mapping(source = "createdBy", target = "createdBy")
    @Mapping(target = "createdTime", source = "createdDate")
    @Mapping(target = "updatedBy" , source = "lastModifiedBy")
    @Mapping(target = "updatedTime", source = "lastModifiedDate")
    @Mapping(target = "addressInfo", source = "address")
    @Mapping(target = "orderInfoList", source = "orders")
    UserInfo toInfo(User user);

    List<UserInfo> toInfoList(List<User> users);
}
