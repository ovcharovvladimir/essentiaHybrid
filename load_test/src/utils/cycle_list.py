"""
Provide a cycle list.
"""


class CycleList(list):
    """
    Cycle list implementation.
    """

    def __getitem__(self, item):
        """
        If item is an int index, then cycle through elements.
        """
        if isinstance(item, int):
            item = item % (self.__len__())

        return super(CycleList, self).__getitem__(item)
