package com.webank.taskman.base;

import com.baomidou.mybatisplus.core.metadata.IPage;

import java.util.List;

public interface BaseConverter<D, E> {

    E toEntity(D dto);

    D toDto(E entity);

    List<E> toEntity(List<D> dtoList);

    List <D> toDto(List<E> entityList);

}
