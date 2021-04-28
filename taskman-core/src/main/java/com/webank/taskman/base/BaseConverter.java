package com.webank.taskman.base;

import java.util.List;

public interface BaseConverter<D, E> {

    E convertToEntity(D dto);

    D convertToDto(E entity);

    List<E> convertToEntities(List<D> dtos);

    List<D> convertToDtos(List<E> entities);

}
