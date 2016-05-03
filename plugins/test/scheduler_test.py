#!/usr/bin/python -O 


def catalan_expr(tups, ops):
    if len(tups)==1:
        yield tups[0]
    else:
        for op in ops:
            for i in range(1, len(tups)):
                for x in catalan_expr(tups[0:i],ops):
                    for y in catalan_expr(tups[i:], ops):
                        yield '(' + op.join((x,y)) + ')'
    
print len(list(catalan_expr('AB','*'))),list(catalan_expr('AB','*'))
print "="*20
print len(list(catalan_expr('ABC','*'))),list(catalan_expr('ABC','*'))
print "="*20
print len(list(catalan_expr('ABCD','*'))),list(catalan_expr('ABCD','*'))
print "="*20
print "="*20
        
print len(list(catalan_expr('AB','/*'))),list(catalan_expr('AB','/*'))
print "="*20
print len(list(catalan_expr('ABC','/*'))),list(catalan_expr('ABC','/*'))
print "="*20
print len(list(catalan_expr('ABCD','/*'))),list(catalan_expr('ABCD','/*'))

